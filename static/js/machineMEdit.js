$(() => {
    $('#form').on('submit', function (ev) {
        ev.preventDefault();

        const id = $('#machineId').val();
        const data = getObjFromSerialize($(this).serializeArray(), ['id']);
        $.ajax({
            url: '/machineM/save',
            method: parseInt(id) === 0? 'post': 'put',
            dataType: 'json',
            data: JSON.stringify(data)
        }).done((r) => {
            if(r && r.id){
                window.location.assign(`/machineM/save?id=${r.id}`);
            }else{
                BootstrapDialog.alert('No data return.');
            }
        }).fail(alertErr);
    });
});