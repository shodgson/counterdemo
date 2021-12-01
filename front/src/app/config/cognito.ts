const awsconfig = {
  Region: import.meta.env.VITE_AWS_REGION,
  UserPoolId: import.meta.env.VITE_AWS_USER_POOL_ID,
  ClientId: import.meta.env.VITE_AWS_USER_POOL_CLIENT_ID,
  // IdentityPoolId:
};

export default awsconfig;
