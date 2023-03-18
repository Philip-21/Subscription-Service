function selectPlan(x, plan) {
      //sweet alert fires a dialogue box.
      Swal.fire({
          title: 'Subscribe',
          html: 'Are you sure you want to subscribe to the ' + plan + '?',
          showCancelButton: true,
          confirmButtonText: 'Subscribe',
      }).then((result) => {
        //the click action calls the subcribe to plan handler in the backend
          if (result.isConfirmed) {
              window.location.href = '/members/subscribe?id=' + x;
          }
      })
  }