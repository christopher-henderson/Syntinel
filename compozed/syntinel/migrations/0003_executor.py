# -*- coding: utf-8 -*-
# Generated by Django 1.10.3 on 2017-02-21 19:24
from __future__ import unicode_literals

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('syntinel', '0002_auto_20170215_2223'),
    ]

    operations = [
        migrations.CreateModel(
            name='Executor',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('hostName', models.CharField(max_length=256)),
                ('port', models.CharField(max_length=5)),
                ('Scheme', models.CharField(default='http', max_length=5)),
            ],
        ),
    ]
