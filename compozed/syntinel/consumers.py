# In consumers.py
from channels import Group
from channels.sessions import channel_session

# Connected to websocket.connect
@channel_session
def ws_connect(message):
    # Accept connection
    message.reply_channel.send({"accept": True})
    # Work out room name from path (ignore slashes)
    # /testRun/console/id
    # /testRun/console/id-executor
    room = message.content['path'].split("/")[-1]
    # Save room in session and add us to the group
    if "executor" in room:
        testID = room.split("-")[0]
        message.channel_session['room'] = testID
        print("Executor in room: " + testID)
    else:
        print("Consumer in room: " + room)
        message.channel_session['room'] = room
        Group(room).add(message.reply_channel)

# Connected to websocket.receive
@channel_session
def ws_message(message):
    room = message.channel_session['room']
    Group(room).send({
        "text": message["text"]
    })

# Connected to websocket.disconnect
@channel_session
def ws_disconnect(message):
    Group("chat-%s" % message.channel_session['room']).discard(message.reply_channel)
