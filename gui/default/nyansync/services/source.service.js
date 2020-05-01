(function () {
    'use strict';

    angular.module('app')
        .factory('SourceService', ['$resource', 'urls', function ($resource, urls) {
            return $resource(urls.BASE_API + 'source/:id', {}, {
                query: {
                    method: 'GET',
                    params: { id: '' },
                    isArray: true,
                },
            });
        }]);

})();
