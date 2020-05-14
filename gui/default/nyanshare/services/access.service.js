(function() {
  'use strict';

  angular.module('app')
    .factory('AccessService', ['$resource', 'urls', function( $resource, urls ) {
      return $resource(urls.BASE_API + 'access/:Id', {Id: '@_orig_id'}, {
        query: {
          method: 'GET',
          params: { Id: '' },
          isArray: true,
          cache : true,
        },
      });
    }]);

})();
