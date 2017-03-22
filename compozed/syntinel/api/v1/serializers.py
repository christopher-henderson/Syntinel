from rest_framework.serializers import ModelSerializer, CharField, PrimaryKeyRelatedField
from syntinel.models import Test, Project, TestRun, Executor


class TestSerializer(ModelSerializer):

    class Meta:
        model = Test
        fields = '__all__'


class ProjectSerializer(ModelSerializer):

    tests = PrimaryKeyRelatedField(many=True, read_only=True)

    class Meta:
        model = Project
        fields = '__all__'


class TestRunSerializer(ModelSerializer):

    class Meta:
        model = TestRun
        fields = '__all__'


class ExecutorSerializer(ModelSerializer):

    class Meta:
        model = Executor
        fields = '__all__'
