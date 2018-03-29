let ruleData = [];
let exampleMap = {};
$(() => {
   $('#form').on('submit', function(ev){
        ev.preventDefault();

        let data = getObjFromSerialize($(this).serializeArray(), ['quantity', 'outputCount', 'idRuleId']);

        $.ajax({
            url: '/jobM/save',
            method: 'post',
            dataType: 'json',
            data: JSON.stringify(data)
        }).done((r) => {
            const data = r.data || {};
            if(data.id){
                window.location.assign(`/jobM/save?id=${data.id}`)
            }else{
                BootstrapDialog.alert('No Id in return message.');
            }
        }).fail(alertErr)
   });

   if(!window.location.search){
        $.ajax({
            url: '/api/v1/jobIdRuleList',
        }).done((r) => {
            ruleData = r.data || [];
            exampleMap = ruleData.reduce((pre, cur) => {
                pre[cur.id] = cur.example;
                return pre;
            }, {});
        }).fail(alertErr);
   }
});

function showExample(){
    const ruleId = $('#idRuleSelect').val();
    const exampleSpan = $('#example');
    if(ruleId){
        exampleSpan.html('<br>example: ' + exampleMap[ruleId]);
    }else{
        exampleSpan.html('');
    }
}
