package sync

import (
	"sync"
	"time"

	"github.com/icon-project/goloop/common/log"
	"github.com/icon-project/goloop/module"
)

type client struct {
	ph    module.ProtocolHandler
	mutex sync.Mutex
	log   log.Logger
}

func (cl *client) hasNode(p *peer, wsHash, prHash, nrHash, vh []byte,
	expiredCb func(pi module.ProtocolInfo, b []byte, p *peer)) error {
	reqID := p.reqID + 1
	msg := &hasNode{reqID, wsHash, vh, prHash, nrHash}
	b, _ := c.MarshalToBytes(msg)
	if err := cl.ph.Unicast(protoHasNode, b, p.id); err != nil {
		cl.log.Infof("Failed to request for protoHasNode err(%+v)\n", err)
		return err
	}
	p.reqID = reqID
	p.timer = time.AfterFunc(time.Millisecond*configExpiredTime, func() {
		r := &result{p.reqID, ErrTimeExpired}
		b, _ := c.MarshalToBytes(r)
		cl.log.Debugf("hasNode time expired for p(%s)\n", p)
		expiredCb(protoResult, b, p)
	})
	return nil
}

func (cl *client) requestNodeData(p *peer, hash [][]byte, t syncType,
	expiredCb func(pi module.ProtocolInfo, b []byte, p *peer)) error {
	reqID := p.reqID + 1
	msg := &requestNodeData{reqID, t, hash}
	b, _ := c.MarshalToBytes(msg)
	if err := cl.ph.Unicast(protoRequestNodeData, b, p.id); err != nil {
		cl.log.Infof("Failed to request for protoRequestNodeData err(%+v)\n", err)
		return err
	}

	p.reqID = reqID
	p.timer = time.AfterFunc(time.Millisecond*configExpiredTime, func() {
		nd := &nodeData{p.reqID, ErrTimeExpired, t, nil}
		b, _ := c.MarshalToBytes(nd)
		cl.log.Debugf("requestNodeData time expired, peer(%s)\n", p)
		expiredCb(protoNodeData, b, p)
	})
	return nil
}

func newClient(ph module.ProtocolHandler, log log.Logger) *client {
	cl := new(client)
	cl.ph = ph
	cl.log = log
	return cl
}