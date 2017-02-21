from rest_framework.serializers import ModelSerializer, CharField
from syntinel.models import Test, Suite, TestRun, Executor


class TestSerializer(ModelSerializer):

    class Meta:
        model = Test
        fields = '__all__'


class SuiteSerializer(ModelSerializer):

    class Meta:
        model = Suite
        fields = '__all__'


class TestRunSerializer(ModelSerializer):

    class Meta:
        model = TestRun
        fields = '__all__'


class ExecutorSerializer(ModelSerializer):

    class Meta:
        model = Executor
        fields = '__all__'
