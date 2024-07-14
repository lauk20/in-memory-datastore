import key_val_pb2
import key_val_pb2_grpc
import pubsub_pb2
import pubsub_pb2_grpc
import grpc

def create_key_val_pair(key, value):
    k = key_val_pb2.Key(val=key)
    v = key_val_pb2.Value(val=value, found=True)
    return key_val_pb2.KeyValuePair(key=k, value=v)

def create_pub_sub_message(topic, message):
    return pubsub_pb2.Pub(topic=topic, msg=message)

def publish(stub, topic, message):
    return sub.Publish(create_pub_sub_message(topic, message))

def subscribe(stub, topic):
    for submessage in stub.Subscribe(pubsub_pb2.String(msg=topic)):
        print(submessage.msg)

def run():
    channel = grpc.insecure_channel('localhost:6379')
    stub = key_val_pb2_grpc.KeyValueStub(channel)
    stub.SetValue(create_key_val_pair("New York", "Metro"))
    print(stub.GetValue(key_val_pb2.Key(val="New York")).val)

    pubsubStub = pubsub_pb2_grpc.PubSubStub(channel)
    subscribe(pubsubStub, "test_topic")
    """
    pubsubStub = pubsub_pb2_grpc.PubSubStub(channel)
    print(publish(pubsubStub, "test_topic", "secret message through pubsub"), "subscribers")
    """

if __name__ == "__main__":
    run()