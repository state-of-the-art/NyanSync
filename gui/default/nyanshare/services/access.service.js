(function() {
  'use strict';

  angular.module('app')
    .factory('AccessService', ['$resource', 'urls', function( $resource, urls ) {
      return $resource(urls.BASE_API + 'access/:Id', {Id: '@Id'}, {
        query: {
          method: 'GET',
          params: { Id: '' },
          isArray: true,
          cache : true,
        },
      });
    }]);

})();
