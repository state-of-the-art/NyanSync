(function() {
  'use strict';

  angular
    .module('app')
    .controller('RoleController', ['role', '$scope', '$uibModalInstance', '$uibModal', 'RoleService',
      function( role, $scope, $uibModalInstance, $uibModal, RoleService ) {
        var vm = this;

        if( role instanceof RoleService )
          vm.title = 'Edit role "' + role.Id + '"';
        else
          vm.title = 'Create new role';

        (function initController() {
          // Role to create or edit
          $scope.role = role instanceof RoleService ? role : new RoleService(role);

          // Set original id for renaming
          $scope.role._orig_id = $scope.role.Id;

          // To validate the used role ids
          $scope.used_ids = []
          RoleService.query().$promise.then(function(roles){
            // Make sure while edit we can reuse the same id of role
            $scope.used_ids = roles.map(r => r.Id != $scope.role.Id ? r.Id : null);
          });
        })();

        vm.submit = function() {
          if( !$scope.role._orig_id )
            $scope.role._orig_id = $scope.role.Id

          $scope.role.$save().then(function() {
            $uibModalInstance.close($scope.role);
          });
        };

        vm.cancel = function() {
          $uibModalInstance.dismiss('cancel');
        };

        vm.remove = function() {
          $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: 'nyanshare/modal/modal.question.html',
            controller: 'QuestionController',
            controllerAs: 'vm',
            size: 'sm',
            resolve: {
              body: function(){
                return 'Are you sure you want to remove role "' + $scope.role.Id + '"?';
              },
            },
          }).result.then(function( result ) {
            if( result === true ) {
              $scope.role.$remove().then(function(){
                $uibModalInstance.close();
              });
            }
          });
        };
      }
    ])
    .directive('myRoleUnique', function() {
      return {
        require: 'ngModel',
        scope: true,
        link: function(scope, elm, attrs, ctrl) {
          if( !ctrl )
            return;
          ctrl.$validators.roleunique = function(modelValue, viewValue) {
            return scope.used_ids.indexOf(modelValue) === -1;
          };
        },
      };
    });

})();
