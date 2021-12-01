const awsconfig = {
    Auth: {

        // REQUIRED only for Federated Authentication - Amazon Cognito Identity Pool ID
        //identityPoolId: 'XX-XXXX-X:XXXXXXXX-XXXX-1234-abcd-1234567890ab',
        
        // REQUIRED - Amazon Cognito Region
        //region: 'us-east-2',

        // OPTIONAL - Amazon Cognito Federated Identity Pool Region 
        // Required only if it's different from Amazon Cognito Region
        //identityPoolRegion: 'us-east-2',

        // OPTIONAL - Amazon Cognito User Pool ID
        //userPoolId: 'us-east-2_59jZflOMc',

        // OPTIONAL - Amazon Cognito Web Client ID (26-char alphanumeric string)
        //userPoolWebClientId: '1t8hrv56fgkeagfskub8uhvvrm',

        // OPTIONAL - Enforce user authentication prior to accessing AWS resources or not
        //mandatorySignIn: false,

        // OPTIONAL - Configuration for cookie storage
        // Note: if the secure flag is set to true, then the cookie transmission requires a secure protocol
        //cookieStorage: {
        // REQUIRED - Cookie domain (only required if cookieStorage is provided)
            //domain: '.yourdomain.com',
        // OPTIONAL - Cookie path
            //path: '/',
        // OPTIONAL - Cookie expiration in days
            //expires: 365,
        // OPTIONAL - See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie/SameSite
            //sameSite: "strict",
        // OPTIONAL - Cookie secure flag
        // Either true or false, indicating if the cookie transmission requires a secure protocol (https).
            //secure: true
        //},

        // OPTIONAL - customized storage object
        //storage: MyStorage,
        
        // OPTIONAL - Manually set the authentication flow type. Default is 'USER_SRP_AUTH'
        //authenticationFlowType: 'USER_PASSWORD_AUTH',

        // OPTIONAL - Manually set key value pairs that can be passed to Cognito Lambda Triggers
        //clientMetadata: { myCustomKey: 'myCustomValue' },

         // OPTIONAL - Hosted UI configuration
        //oauth: {
            //domain: 'your_cognito_domain',
            //scope: ['phone', 'email', 'profile', 'openid', 'aws.cognito.signin.user.admin'],
            //redirectSignIn: 'http://localhost:3000/',
            //redirectSignOut: 'http://localhost:3000/',
            //responseType: 'code' // or 'token', note that REFRESH token will only be generated when the responseType is code
        //}
    //}
//}

/*
const awsconfig = {
  ClientId: '4aqa7kavaig216g45l94017tt5',
  AppWebDomain : 'linkydink-dev.auth.us-east-2.amazoncognito.com',
  TokenScopesArray : ['email', 'profile','openid'],
  RedirectUriSignIn : 'http://localhost:3000/console/signin',
  RedirectUriSignOut : 'http://localhost:3000/console/signout',
  IdentityProvider : 'COGNITO', // e.g. 'Facebook',
  UserPoolId : 'us-east-2_wvtbM3dh0', // Your user pool id here
  AdvancedSecurityDataCollectionFlag : false,
  //Storage: '<TODO the storage object>',
}
*/

export default awsconfig;
