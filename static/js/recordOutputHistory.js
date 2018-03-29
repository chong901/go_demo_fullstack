let socket = io();

socket.on('reportOutputUpdate', function (dataStr) {
    const data = JSON.parse(dataStr);
    if (data.id) {
        let tr = $('*[data-recordid="' + data.id + '"]').closest('tr');
        if(tr && tr.length !== 0){
            setHtml(tr, 'qty', data.quantity);
            setHtml(tr, 'editor', data.user);
        }
    }
});

function editQty(dom, id) {
    BootstrapDialog.show({
        closable: false,
        title: $('#strEditQty').text(),
        message: createUpdateForm(id),
        onshown: function (bd) {
            $("#editQtyForm").on('submit', function (ev) {
                ev.preventDefault();
                let data = getObjFromSerialize($('#editQtyForm').serializeArray(), ['quantity']);
                $.ajax({
                    url: '/reportOutput?id='+id,
                    method: 'PUT',
                    dataType: 'json',
                    data: JSON.stringify(data)
                }).done(function () {
                    bd.close();
                }).fail(alertErr);
            });
        },
        buttons: [
            {
                label: $('#strBtnConfirm').text(),
                action: function () {
                    $('#editQtyFormSubmit').click();
                }
            },
            {
                label: $('#strBtnCancel').text(),
                action: function (dialogItself) {
                    dialogItself.close();
                }
            }]
    });
}

function createUpdateForm(id) {
    let html = `<form id="editQtyForm">`;
    html += `<label for="quantity">${$('#strQuantity').text()}</label>`;
    html += `<input id="quantity" class="form-control" type="number" name="quantity" required>`;
    html += `<input id="editQtyFormSubmit" type="submit" style="display: none">`;
    html += '</form>';
    return html;
}
