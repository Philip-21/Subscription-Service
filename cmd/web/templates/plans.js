function selectPlan(x, plan) {
      Swal.fire({
          title: 'Subscribe',
          html: 'Are you sure you want to subscribe to the ' + plan + '?',
          showCancelButton: true,
          confirmButtonText: 'Subscribe',
      }).then((result) => {
          if (result.isConfirmed) {
              window.location.href = '/subscribe?id=' + x;
          }
      })
  }