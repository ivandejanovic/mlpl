(function(){
	$(document).ready(function(){
        var languages = [{locale:'sr', label:'srpski'},{locale:'en', label:'english'}];
        var selectedLanguage = 0;

        $('#selectedLang').html(languages[selectedLanguage].label);

        $('#languagePicker > li > a').click(function(){
                var targetLang = $(this).attr('data-lang'),
                    targetDivID = '#content_' + targetLang,
                    parentLi = $(this).parent('li');

                if (!parentLi.hasClass('active')){
                    selectedLanguage = languages.indexOf(targetDivID);
                    parentLi.addClass('active').siblings().removeClass('active');
                    $(targetDivID).toggleClass('in').siblings().removeClass('in');
                }
            }
        );

    });
}())