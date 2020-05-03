(function () {
    'use strict';

    angular.module('app', [
        'ngStorage',
        'ngRoute',
        'ngResource',
        'ui.bootstrap',
    ])
		.constant('urls', {
            BASE: '/',
            BASE_API: '/api/v1/'
        })
        .config(['$routeProvider', '$httpProvider', '$resourceProvider', '$qProvider',
            function($routeProvider, $httpProvider, $resourceProvider, $qProvider) {
            $routeProvider
                .when('/', {
                    controller: 'HomeController',
                    templateUrl: 'nyansync/home/home.view.html',
                    controllerAs: 'vm',
                })
                .when('/login', {
                    controller: 'LoginController',
                    templateUrl: 'nyansync/login/login.view.html',
                    controllerAs: 'vm',
                })
                .otherwise({ redirectTo: '/' });

            $httpProvider.interceptors.push(['$q', '$location', '$localStorage', 'FlashService', '$cacheFactory',
                function ($q, $location, $localStorage, FlashService, $cacheFactory) {
                    return {
                        'request': function (config) {
                            // Clean cache if we need to
                            if( config.cache && config.params && config.params.cache === false ) {
                                $cacheFactory.get('$http').remove(config.url);
                                delete config.params.cache;
                            }

                            // Inject auth header
                            config.headers = config.headers || {};
                            if ($localStorage.token) {
                                config.headers.Authorization = 'Bearer ' + $localStorage.token;
                            }
                            return config;
                        },
                        'response': function (res) {
                            if (res.data.data)
                                res.data = res.data.data;
                            return res;
                        },
                        'responseError': function (res) {
                            FlashService.Error(res.data.message);
                            if (res.status === 401 || res.status === 403) {
                                delete $localStorage.token;
                                $location.path('/login');
                            }
                            return $q.reject(res);
                        },
                    };
                }]);

            $resourceProvider.defaults.stripTrailingSlashes = false;
            $qProvider.errorOnUnhandledRejections(false);
        }])
        .run(function($rootScope, $location, $localStorage) {
            $rootScope.$on( "$routeChangeStart", function(event, next) {
                if ($localStorage.token == null) {
                    $location.path("/login");
                }
            });
        });

})();
