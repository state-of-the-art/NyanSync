(function() {
  'use strict';

  angular
    .module('app')
    .controller('UserController', ['user', '$scope', 'UserService', '$uibModalInstance', '$uibModal',
      function( user, $scope, UserService, $uibModalInstance, $uibModal ) {
        var vm = this;

        if( user && user.Login )
          vm.title = 'Edit user "' + user.Name + '"';
        else
          vm.title = 'Create new user';

        (function initController() {
          // User to create or edit
          $scope.user = user instanceof UserService ? user : new UserService(user);
          $scope.user._orig_login = $scope.user.Login;
          if( !$scope.user._orig_login )
            $scope.user.Login = $scope.user.Name;

          // To validate the used user logins
          UserService.query().$promise.then(function(users){
            // Make sure while edit we can reuse the same login of user
            $scope.used_logins = users.map(u => u.Login != $scope.user.Login ? u.Login : null);
          });
        })();

        vm.submit = function() {
          if( !$scope.user._orig_login )
            $scope.user._orig_login = $scope.user.Login
          $scope.user.$save().then(function() {
            $uibModalInstance.close($scope.user);
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
                return 'Are you sure you want to remove user "' + $scope.user.Name +
                  '" with login "' + $scope.user._orig_login + '"?';
              },
            },
          }).result.then(function( result ) {
            if( result === true ) {
              $scope.user.$remove().then(function(){
                $uibModalInstance.close();
              });
            }
          });
        };
      }
    ])
    .directive('myUserUnique', function() {
      return {
        require: 'ngModel',
        scope: true,
        link: function(scope, elm, attrs, ctrl) {
          if( !ctrl )
            return;
          ctrl.$validators.userunique = function(modelValue, viewValue) {
            return scope.used_logins.indexOf(modelValue) === -1;
          };
        },
      };
    });

})();
