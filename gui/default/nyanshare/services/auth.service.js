(function() {
  'use strict';

  angular.module('app')
    .factory('AuthService', ['$http', '$rootScope', '$localStorage', '$location', 'urls',
      function( $http, $rootScope, $localStorage, $location, urls ) {
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
          if( $localStorage.account ) {
            var token = $localStorage.account.token;
            var encoded = token.split('.')[1];
            var claims = JSON.parse(urlBase64Decode(encoded));
            return claims;
          }
          return {};
        }

        var tokenClaims = getClaimsFromToken();

        var logout = function( notify_server ) {
          if( notify_server )
            $http.post(urls.BASE_API + 'auth/logout');
          tokenClaims = {};
          delete $localStorage.account;
          delete $rootScope.account;
          $location.path("/login");
        };

        var refreshRootScope = function() {
          if( $localStorage.account ) {
            $localStorage.account.Logout = logout;
            $rootScope.account = $localStorage.account;
          }
        };

        return {
          Login: function( login, password, success, error ) {
            $http.post(urls.BASE_API + 'auth/login', {
              "login": login,
              "password": password,
            }).then(function(res) {
              $localStorage.account = res.data.user;
              $localStorage.account.token = res.data.token;
              tokenClaims = getClaimsFromToken();
              refreshRootScope();
              success();
            }, error);
          },
          Logout: logout,
          GetTokenClaims: function() {
            return tokenClaims;
          },
          RefreshRootScope: refreshRootScope,
        };
      }
    ]);

})();
