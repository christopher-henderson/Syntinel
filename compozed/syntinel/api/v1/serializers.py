from rest_framework.serializers import ModelSerializer, CharField
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

    content = CharField(max_length=2 ** 56)

    class Meta:
        model = Script
        fields = '__all__'

    def create(self, validated_data):
        content = validated_data.get('content')
        del validated_data['content']
        script = self.Meta.model.objects.create(**validated_data)
        script.content = content
        return script

    def update(self, script, validated_data):
        content = validated_data.get('content')
        script.content = content
        return script


class TestRunSerializer(ModelSerializer):

    class Meta:
        model = TestRun
        fields = '__all__'
