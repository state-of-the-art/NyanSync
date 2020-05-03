(function () {
    'use strict';

    angular.module('app')
        .factory('SourceService', ['$resource', 'urls', function ($resource, urls) {
            return $resource(urls.BASE_API + 'source/:Id', {Id: '@Id'}, {
                query: {
                    method: 'GET',
                    params: { Id: '' },
                    isArray: true,
                    cache : true,
                },
            });
        }]);

})();
