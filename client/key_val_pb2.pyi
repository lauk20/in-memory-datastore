from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Key(_message.Message):
    __slots__ = ("val",)
    VAL_FIELD_NUMBER: _ClassVar[int]
    val: str
    def __init__(self, val: _Optional[str] = ...) -> None: ...

class Value(_message.Message):
    __slots__ = ("val", "found")
    VAL_FIELD_NUMBER: _ClassVar[int]
    FOUND_FIELD_NUMBER: _ClassVar[int]
    val: str
    found: bool
    def __init__(self, val: _Optional[str] = ..., found: bool = ...) -> None: ...

class KeyValuePair(_message.Message):
    __slots__ = ("key", "value")
    KEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    key: Key
    value: Value
    def __init__(self, key: _Optional[_Union[Key, _Mapping]] = ..., value: _Optional[_Union[Value, _Mapping]] = ...) -> None: ...
