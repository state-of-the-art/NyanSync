(function () {
  'use strict';

  angular
    .module('app')
    .controller('QuestionController', ['body', '$uibModalInstance',
      function (body, $uibModalInstance) {
          var vm = this;

          vm.body = body

          vm.confirm = function () {
            $uibModalInstance.close(true);
          };

          vm.cancel = function () {
            $uibModalInstance.dismiss('cancel');
          };
        }
      ]);

})();
