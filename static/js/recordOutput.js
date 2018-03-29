let socket = io();

socket.on('reportOutputUpdate', function (dataStr) {
    const data = JSON.parse(dataStr);
    if (data.jobId) {
        let tr = $('*[data-jobid="' + data.jobId + '"]').closest('tr');
        if(tr && tr.length !== 0){
            let oldQty = parseInt(tr.find('.reportedOutput').text());
            let updateQty = oldQty + parseInt(data.quantity);
            setHtml(tr, 'reportedOutput', updateQty.toString());
        }
    }
});

function reportOutput(dom, id) {
    BootstrapDialog.show({
        closable: false,
        title: $('#strReportOutput').text(),
        message: createUpdateForm(id),
        onshown: function (bd) {
            $("#reportOutputForm").on('submit', function (ev) {
                ev.preventDefault();
                let data = getObjFromSerialize($('#reportOutputForm').serializeArray(), ['quantity']);
                $.ajax({
                    url: '/reportOutput',
                    method: 'POST',
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
                    $('#reportOutputFormSubmit').click();
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
    let html = `<form id="reportOutputForm">`;
    html += `<label for="quantity">${$('#strQuantity').text()}</label>`;
    html += `<input id="quantity" class="form-control" type="number" name="quantity" required>`;
    html += `<input type="hidden" name="jobId" value="${id}">`;
    html += `<input id="reportOutputFormSubmit" type="submit" style="display: none">`;
    html += '</form>';
    return html;
}
