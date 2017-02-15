# -*- coding: utf-8 -*-
# Generated by Django 1.10.3 on 2017-02-15 22:19
from __future__ import unicode_literals

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Suite',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=65536)),
            ],
        ),
        migrations.CreateModel(
            name='Test',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=65536)),
                ('script', models.CharField(max_length=65536)),
                ('dockerfile', models.CharField(max_length=65536)),
                ('environmentVariables', models.CharField(max_length=65536, null=True)),
                ('health', models.IntegerField(default=100)),
                ('suite', models.ForeignKey(null=True, on_delete=django.db.models.deletion.CASCADE, to='syntinel.Suite')),
            ],
        ),
        migrations.CreateModel(
            name='TestRun',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('log', models.CharField(max_length=65536)),
                ('successful', models.BooleanField()),
                ('test', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='syntinel.Test')),
            ],
        ),
    ]
