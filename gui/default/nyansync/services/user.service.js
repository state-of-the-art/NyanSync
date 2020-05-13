(function () {
    'use strict';

    angular.module('app')
        .factory('UserService', ['$resource', 'urls', function ($resource, urls) {
            return $resource(urls.BASE_API + 'user/:Id', {Id: '@_orig_id'}, {
                query: {
                    method: 'GET',
                    params: { Id: '' },
                    isArray: true,
                    cache : true,
                },
            });
        }]);

})();
