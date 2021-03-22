import * as cdk from '@aws-cdk/core';
import * as apigateway from '@aws-cdk/aws-apigateway';
import * as lambda from '@aws-cdk/aws-lambda';

export class LotteryStack extends cdk.Stack {
  constructor(scope: cdk.App, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const lotteryHandler = new lambda.Function(this, 'SlackLotteryBot-LotteryFunction', {
      runtime: lambda.Runtime.GO_1_X,
      handler: 'SlackLotteryBot-LotteryHandler',
      code: lambda.Code.fromAsset('./dist/lottery.zip'),
      environment: {
        'BOTTOKEN': process.env.BOTTOKEN || '',
        'VERIFICATIONTOKEN': process.env.VERIFICATIONTOKEN || '',
        'OAUTHTOKEN': process.env.OAUTHTOKEN || '',
      }
    })

    const selectHandler = new lambda.Function(this, 'SlackLotteryBot-SelectFunction', {
      runtime: lambda.Runtime.GO_1_X,
      handler: 'SlackLotteryBot-SelectHandler',
      code: lambda.Code.fromAsset('./dist/select.zip'),
      environment: {
        'BOTTOKEN': process.env.BOTTOKEN || '',
        'VERIFICATIONTOKEN': process.env.VERIFICATIONTOKEN || '',
        'OAUTHTOKEN': process.env.OAUTHTOKEN || '',
      }
    })

    new apigateway.LambdaRestApi(this, 'SlackLotteryBot-LotteryGateway', {
      handler: lotteryHandler,
      restApiName: 'SlackLotteryBot-LotteryAPI'
    })

    new apigateway.LambdaRestApi(this, 'SlackLotteryBot-SelectGateway', {
      handler: selectHandler,
      restApiName: 'SlackLotteryBot-SelectAPI'
    })
  }
}
