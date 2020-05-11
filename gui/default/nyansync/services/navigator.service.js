(function () {
  'use strict';

  angular.module('app')
    .factory('NavigatorService', ['$resource', 'urls', function ($resource, urls) {
      return $resource(urls.BASE_API + 'navigate/:path', {path: '@path'}, {
        query: { // Input params - array with path items
          method: 'GET',
          isArray: true,
          cache : false, // TODO: enable to optimze
          interceptor: {
            request: function(config) {
              if( config.params ) {
                config.url += Object.values(config.params).join('/');
                delete config.params;
              }
              return config;
            },
          },
        },
      });
    }]);

})();
