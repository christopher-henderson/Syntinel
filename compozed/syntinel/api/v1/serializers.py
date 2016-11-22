from rest_framework.serializers import ModelSerializer
from .models import Docker, Test, Suite, Script, TestRun


class DockerSerializer(ModelSerializer):

    class Meta:
        model = Docker
        fields = '__all__'


class TestSerializer(ModelSerializer):

    class Meta:
        model = Test
        fields = '__all__'


class SuiteSerializer(ModelSerializer):

    class Meta:
        model = Suite
        fields = '__all__'


class ScriptSerializer(ModelSerializer):

    class Meta:
        model = Script
        fields = '__all__'


class TestRunSerializer(ModelSerializer):

    class Meta:
        model = TestRun
        fields = '__all__'
