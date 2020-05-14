(function() {
  'use strict';

  angular.module('app')
    .factory('AuthService', ['$http', '$localStorage', 'urls',
      function( $http, $localStorage, urls ) {
        function urlBase64Decode(str) {
          var output = str.replace('-', '+').replace('_', '/');
          switch( output.length % 4 ) {
            case 0:
              break;
            case 2:
              output += '==';
              break;
            case 3:
              output += '=';
              break;
            default:
              throw 'Illegal base64url string!';
          }
          return window.atob(output);
        }

        function getClaimsFromToken() {
          var token = $localStorage.token;
          var user = {};
          if( typeof token !== 'undefined' ) {
            var encoded = token.split('.')[1];
            user = JSON.parse(urlBase64Decode(encoded));
          }
          return user;
        }

        var tokenClaims = getClaimsFromToken();

        return {
          Login: function( login, password, success, error ) {
            $http.post(urls.BASE_API + 'auth/login', {
              "login": login,
              "password": password,
            }).then(success, error);
          },
          Logout: function( success, error ) {
            $http.post(urls.BASE_API + 'auth/logout').then(success, error)
            tokenClaims = {};
            delete $localStorage.token;
            success();
          },
          GetTokenClaims: function() {
            return tokenClaims;
          }
        };
      }
    ]);

})();
