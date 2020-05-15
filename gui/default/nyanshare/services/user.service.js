(function() {
  'use strict';

  angular.module('app')
    .factory('UserService', ['$resource', 'urls', function( $resource, urls ) {
      return $resource(urls.BASE_API + 'user/:Login', {Login: '@_orig_login'}, {
        query: {
          method: 'GET',
          params: { Login: '' },
          isArray: true,
          cache : true,
        },
      });
    }]);

})();
