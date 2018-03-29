$(() => {
    $('#form').on('submit', function (ev) {
        ev.preventDefault();
        let data = getObjFromSerialize($(this).serializeArray(), ['id', 'category', 'digitNumber'], null, ['isDateShow']);
        const id = parseInt($('#id').val());
        $.ajax({
            url: '/idRule',
            method: id? 'put':'post',
            dataType: 'json',
            data: JSON.stringify(data)
        }).done((r) => {
            const data = r.data || {};
            window.location.assign(`${window.location.pathname}?id=${data.id}`)
        }).fail(alertErr);
    });

    showExample();
});

function showExample(){
    const prefix = $('input[name="prefix"]').val();
    const digitNumber = parseInt($('input[name="digitNumber"]').val());
    const isDate = parseInt($('input[name="isDateShow"]:checked').val());

    let example = prefix;
    if(isDate){
        example += moment(new Date()).format('YYYYMMDD');
    }
    for (let i = 0; i < digitNumber - 1; i++){
        example += '0';
    }

    example += '1';
    $('#example').html(example);
}