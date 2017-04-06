from rest_framework.serializers import ModelSerializer, CharField, PrimaryKeyRelatedField, ValidationError
from syntinel.models import Test, Project, TestRun, Executor


class TestSerializer(ModelSerializer):

    class Meta:
        model = Test
        fields = '__all__'

    def validate_interval(self, interval):
        if interval is not None and interval < 1:
            raise ValidationError("Interval must be 1 second or greater.")
        return interval


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
