(function () {
    'use strict';

    angular.module('app', [
        'ngStorage',
        'ngRoute',
        'ngResource',
        'ui.bootstrap',
        'ui-notification',
        'ngTagsInput',
    ])
        .constant('urls', {
            BASE: '/',
            BASE_API: '/api/v1/'
        })
        .config(['$routeProvider', '$httpProvider', '$resourceProvider', '$qProvider', 'NotificationProvider',
            function($routeProvider, $httpProvider, $resourceProvider, $qProvider, NotificationProvider) {
            NotificationProvider.setOptions({
                delay: 10000,
                startTop: 20,
                startRight: 10,
                verticalSpacing: 20,
                horizontalSpacing: 20,
                positionX: 'right',
                positionY: 'top'
            });

            $routeProvider
                .when('/', {
                    controller: 'HomeController',
                    templateUrl: 'nyanshare/home/home.view.html',
                    controllerAs: 'hm',
                })
                .when('/login', {
                    controller: 'LoginController',
                    templateUrl: 'nyanshare/login/login.view.html',
                    controllerAs: 'vm',
                })
                .otherwise({ redirectTo: '/' });

            $httpProvider.interceptors.push(['$q', '$injector', '$location', '$localStorage', '$cacheFactory',
                function ($q, $injector, $location, $localStorage, $cacheFactory) {
                    return {
                        'request': function (config) {
                            // Clean cache if we need to
                            if( config.cache && config.params && config.params.cache === false ) {
                                $cacheFactory.get('$http').remove(config.url);
                                delete config.params.cache;
                            }

                            // Inject auth header
                            config.headers = config.headers || {};
                            if( $localStorage.token && config.url.startsWith('/api') ) {
                                config.headers.Authorization = 'Bearer ' + $localStorage.token;
                            }
                            return config;
                        },
                        'response': function (res) {
                            // TODO: Do not show when cache is used
                            if (res.data.message)
                                $injector.get('Notification').success('API: ' + res.data.message);
                            if (res.data.data)
                                res.data = res.data.data;
                            return res;
                        },
                        'responseError': function (res) {
                            $injector.get('Notification').error('API: ' + res.data.message);
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
