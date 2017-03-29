# In consumers.py
from channels import Group
from channels.sessions import channel_session
from django.core.cache import cache
from syntinel.models import TestRun

import logging

logger = logging.getLogger("consumer")

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
        LogCache.getLogCache(testID)
        message.channel_session['room'] = testID
        logger.debug("Executor in room: " + testID)
    else:
        logger.debug("Consumer in room: " + room)
        message.channel_session['room'] = room
        logCache = cache.get(LogCache.getCacheName(room))
        if logCache is not None:
            message.reply_channel.send({"text": logCache.retrieve()})
        Group(room).add(message.reply_channel)

# Connected to websocket.receive
@channel_session
def ws_message(message):
    room = message.channel_session['room']
    text = message["text"]
    logCache = LogCache.getLogCache(room)
    logCache.update(text)
    Group(room).send({
        "text": text
    })

# Connected to websocket.disconnect
@channel_session
def ws_disconnect(message):
    Group("chat-%s" % message.channel_session['room']).discard(message.reply_channel)


class LogCache(object):

    MAX_CACHE = 100

    @staticmethod
    def getLogCache(id):
        name = LogCache.getCacheName(id)
        logCache = cache.get(name)
        if logCache is None:
            newCache = LogCache(id)
            cache.set(name, newCache)
            return newCache
        return logCache

    @staticmethod
    def getCacheName(id):
        return "Cache_{ID}".format(ID=id)

    def __init__(self, test_run_id):
        self.id = test_run_id
        self.test_run = TestRun.objects.get(id=test_run_id)
        self.number_of_lines_cached = 0
        cache.set(self.id, "")

    def retrieve(self):
        return "{SAVED}{CACHED}".format(
            SAVED=self.test_run.log, CACHED=cache.get(self.id))

    def update(self, log_line):
        if self.number_of_lines_cached >= self.MAX_CACHE:
            self.flush()
        cache.set(self.id, "{CACHED}{NEW}\n".format(
            CACHED=cache.get(self.id),
            NEW=log_line))
        self.number_of_lines_cached += 1
        cache.set(self.getCacheName(self.id), self)

    def flush(self):
        self.test_run.log = "{SAVED}{CACHED}".format(
            SAVED=self.test_run.log,
            CACHED=cache.get(self.id))
        self.test_run.save()
        cache.set(self.id, "")
        self.number_of_lines_cached = 0
        cache.set(self.getCacheName(self.id), self)

    def finalize(self):
        self.flush()
        cache.delete(self.id)
        cache.delete(self.getCacheName(self.id))
