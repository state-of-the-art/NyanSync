(function() {
  'use strict';

  angular.module('app')
    .factory('RoleService', ['$resource', 'urls', function( $resource, urls ) {
      return $resource(urls.BASE_API + 'role/:Id', {Id: '@_orig_id'}, {
        query: {
          method: 'GET',
          params: { Id: '' },
          isArray: true,
          cache : true,
        },
      });
    }]);

})();
