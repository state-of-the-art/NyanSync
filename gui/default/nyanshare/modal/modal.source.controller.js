(function () {
    'use strict';

    angular
        .module('app')
        .controller('SourceController', ['title', 'source', '$scope', 'SourceService', '$uibModalInstance', '$uibModal',
            function (title, source, $scope, SourceService, $uibModalInstance, $uibModal) {
                var vm = this;

                vm.title = title;

                (function initController() {
                    // Source to create or edit
                    $scope.source = source ? source : new SourceService();
                    $scope.source._orig_id = $scope.source.Id

                    // To validate the used source ids
                    SourceService.query().$promise.then(function(sources){
                        // Make sure while edit we can reuse the same name of source
                        $scope.used_ids = sources.map(s => s.Id != $scope.source.Id ? s.Id : null);
                    });
                })();

                vm.submit = function () {
                    if( !$scope.source._orig_id )
                        $scope.source._orig_id = $scope.source.Id
                    $scope.source.$save().then(function(){
                        $uibModalInstance.close();
                    });
                };

                vm.cancel = function () {
                    $uibModalInstance.dismiss('cancel');
                };

                vm.remove = function () {
                    $uibModal.open({
                        animation: true,
                        ariaLabelledBy: 'modal-title',
                        ariaDescribedBy: 'modal-body',
                        templateUrl: 'nyanshare/modal/modal.question.html',
                        controller: 'QuestionController',
                        controllerAs: 'vm',
                        size: 'sm',
                        resolve: {
                            body: function(){ return 'Are you sure you want to remove source with name "' + source.Id + '"?'; },
                        },
                    }).result.then(function (result) {
                        if( result === true )
                            $scope.source.$remove().then(function(){
                                $uibModalInstance.close();
                            });
                    });
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
