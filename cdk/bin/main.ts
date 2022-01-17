#!/usr/bin/env node
import * as cdk from '@aws-cdk/core';
import { LotteryStack } from '../lib/lottery-stack';

const app = new cdk.App();
new LotteryStack(app, 'LotteryStack');
