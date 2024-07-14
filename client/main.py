import key_val_pb2
import key_val_pb2_grpc
import grpc

def create_key_val_pair(key, value):
    k = key_val_pb2.Key(val=key)
    v = key_val_pb2.Value(val=value, found=True)
    return key_val_pb2.KeyValuePair(key=k, value=v)

def run():
    channel = grpc.insecure_channel('localhost:6379')
    stub = key_val_pb2_grpc.KeyValueStub(channel)
    stub.SetValue(create_key_val_pair("New York", "Metro"))
    print(stub.GetValue(key_val_pb2.Key(val="New York")).val)

if __name__ == "__main__":
    run()