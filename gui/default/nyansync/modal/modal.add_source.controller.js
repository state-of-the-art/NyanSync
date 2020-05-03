(function () {
    'use strict';

    angular
        .module('app')
        .controller('AddSourceController', ['$scope', 'SourceService', '$uibModalInstance',
            function ($scope, SourceService, $uibModalInstance) {
                var vm = this;

                (function initController() {
                    // To validate the used source ids
                    SourceService.query().$promise.then(function(sources){
                        $scope.used_ids = sources.map(s => s.Id);
                    });

                    // Our new source to create
                    $scope.source = new SourceService();
                })();

                vm.submit = function () {
                    $scope.source.$save().then(function(){
                        $uibModalInstance.close();
                    });
                };

                vm.cancel = function () {
                    $uibModalInstance.dismiss('cancel');
                };
            }
        ])
        .directive('mySourceId', function() {
            return {
                require: 'ngModel',
                scope: true,
                link: function(scope, elm, attrs, ctrl) {
                    if (! ctrl)
                        return;
                    ctrl.$validators.sourceid = function(modelValue, viewValue) {
                        return scope.used_ids.indexOf(modelValue) === -1;
                    };
                },
            };
        });

})();
