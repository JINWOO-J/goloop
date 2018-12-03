#!/usr/bin/env python3
import sys
import os.path
from copy import copy
import os

basedir = os.path.abspath(os.path.join(os.path.dirname(__file__), "../.."))
sys.path.append(basedir)
from pyee.ipc import *


class Address(object):
    def __init__(self, obj):
        if isinstance(obj, bytes):
            if len(obj) < 21:
                raise Exception("IllegalFormat")
            self.__bytes = copy(obj)
        elif isinstance(obj, str):
            if len(obj) < 42:
                raise Exception("IllegalFormat")
            prefix = bytes([obj[:2] == "cx"])
            body = bytes.fromhex(obj[2:])
            self.__bytes = prefix + body
        else:
            raise Exception("IllegalFormat")

    @staticmethod
    def from_str(s: str) -> 'Address':
        if len(s) < 42:
            raise Exception("IllegalFormat")
        prefix = bytes([s[:2] == "cx"])
        body = bytes.fromhex(s[2:])
        return Address(prefix + body)

    def to_bytes(self):
        return copy(self.__bytes)

    def __repr__(self):
        body = self.__bytes[1:].hex()
        if self.__bytes[0] == 0:
            return "hx" + body
        else:
            return "cx" + body


class EECodec(Codec):
    __instance = None

    def __new__(cls, *args, **kwargs):
        if cls.__instance is None:
            cls.__instance = object.__new__(cls, *args, **kwargs)
        return cls.__instance

    def encode(self, obj) -> Tuple[int, bytes]:
        if isinstance(obj, Address):
            return TypeTag.ADDRESS, obj.to_bytes()
        raise Exception

    def decode(self, t: int, b: bytes) -> Any:
        if t == TypeTag.ADDRESS:
            return Address(b)
        else:
            raise Exception("UnknownType:" + t)


class TestEE(object):
    def __init__(self, proxy: 'ServiceManagerProxy'):
        self.__proxy = proxy
        proxy.set_codec(EECodec())
        proxy.set_invoke_handler(self.invoke_handler)

    def invoke_handler(self, code: str, _from: 'Address', _to: 'Address',
                       value: int, limit: int, method: str, params: bytes) -> Tuple[int, int, bytes]:
        print(f'invoke_handler(code={repr(code)},from={_from},to={_to},' +
              f'value={value},limit={limit},method={repr(method)},params={params})')
        self.set_value(b"hello", b"world")
        self.get_value(b'hello')
        self.get_value(b'foo')
        self.get_info()
        self.send_event(["LogEvent(int,str,Address)", 1, "TEST"],
                        [Address.from_str("cx0004444444444444444444444444444444444444")])
        return 0, 10, bytes([])

    def get_value(self, k: bytes) -> Tuple[bool, bytes]:
        ret = self.__proxy.get_value(k)
        print(f"get_value({repr(k)}) -> {repr(ret)}")
        return ret

    def set_value(self, k: bytes, v: bytes):
        print(f"set_value({repr(k)},{repr(v)})")
        return self.__proxy.set_value(k, v)

    def get_info(self) -> Any:
        info = self.__proxy.get_info()
        print(f"get_info() -> {info}")
        return info

    def send_event(self, indexed: List[Any], data: List[Any]):
        print(f"send_event({indexed},{data})")
        self.__proxy.send_event(indexed, data)

    def connect(self, addr: str):
        print(f"connect({addr})")
        self.__proxy.connect(addr)
        self.__proxy.send_version(1, os.getpid(), "python")

    def loop(self):
        self.__proxy.loop()


def main():
    proxy = ServiceManagerProxy()
    ee = TestEE(proxy)
    ee.connect("/tmp/ee.socket")
    ee.loop()


if __name__ == "__main__":
    main()