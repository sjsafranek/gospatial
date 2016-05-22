# -*- coding: utf-8 -*-
# Generated by Django 1.9.4 on 2016-05-05 02:30
from __future__ import unicode_literals

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    dependencies = [
        ('auth', '0007_alter_validators_add_error_messages'),
        ('WallOfDenial', '0001_initial'),
    ]

    operations = [
        migrations.CreateModel(
            name='GeoAPIKey',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('url', models.CharField(max_length=150)),
                ('apikey', models.CharField(max_length=150)),
                ('owner', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='auth.Group')),
            ],
        ),
    ]