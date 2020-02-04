/*
 * Copyright 2020 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package foundation.icon.test.score;

import foundation.icon.icx.IconService;
import foundation.icon.icx.data.Address;
import foundation.icon.icx.transport.jsonrpc.RpcArray;
import foundation.icon.icx.transport.jsonrpc.RpcItem;
import foundation.icon.icx.transport.jsonrpc.RpcObject;
import foundation.icon.icx.transport.jsonrpc.RpcValue;
import foundation.icon.test.common.Constants;
import foundation.icon.test.common.Env;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class ChainScore extends Score {
    private static final int CONFIG_AUDIT = 0x2;
    private static final int CONFIG_DEPLOYER_WHITE_LIST = 0x4;

    public ChainScore(IconService service, Env.Chain chain) {
        super(service, chain, Constants.CHAINSCORE_ADDRESS);
    }

    public int getRevision() throws IOException {
        return call("getRevision", null).asInteger().intValue();
    }

    public int getServiceConfig() throws IOException {
        return call("getServiceConfig", null).asInteger().intValue();
    }

    public static boolean isAuditEnabled(int config) {
        return (config & CONFIG_AUDIT) != 0;
    }

    public boolean isAuditEnabled() throws IOException {
        return isAuditEnabled(this.getServiceConfig());
    }

    public static boolean isDeployerWhiteListEnabled(int config) {
        return (config & CONFIG_DEPLOYER_WHITE_LIST) != 0;
    }

    public boolean isDeployerWhiteListEnabled() throws IOException {
        return isDeployerWhiteListEnabled(this.getServiceConfig());
    }

    public boolean isDeployer(Address address) throws IOException {
        RpcObject params = new RpcObject.Builder()
                .put("address", new RpcValue(address))
                .build();
        return call("isDeployer", params).asBoolean();
    }

    public List<Address> getDeployers() throws IOException {
        List<Address> list = new ArrayList<>();
        RpcArray items = call("getDeployers", null).asArray();
        for (RpcItem item : items) {
            list.add(item.asAddress());
        }
        return list;
    }
}