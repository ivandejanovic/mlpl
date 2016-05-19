/**
 * Created by misa on 18.5.2016.
 */
// Code goes here

(function(){
    $(document).ready(function(){

        //language picker
        var languages = {
                sr: {
                    locale:'sr',
                    langLabel:'srpski',
                    docsLbl:'Dokumentacija',
                    gettingStartedLbl:'Kako početi',
                    flagImgSrc:'images/flag-serbia.jpg'
                },
                en: {
                    locale:'en',
                    langLabel:'english',
                    docsLbl:'Documentation',
                    gettingStartedLbl:'Getting Started',
                    flagImgSrc:'images/flag-united-kingdom.jpg'
                }
            },
            selectedLanguage = 'sr';



        var setNavbarLang = function(selectedLanguage){
            $('.docsLbl').html(languages[selectedLanguage].docsLbl);
            $('.gettingStartedLbl').html(languages[selectedLanguage].gettingStartedLbl);
        };

        var setLanguagePicker = function(){
            $('#selectedLang').html(languages[selectedLanguage].langLabel).prepend($('<img>',{class:'flagImg',src:languages[selectedLanguage].flagImgSrc}));
        };

        setLanguagePicker(selectedLanguage);
        setNavbarLang(selectedLanguage);

        $('#languagePicker > li > a').click(function(){
                var targetLang = $(this).attr('data-lang'),//sr or en
                    targetDivID = '#content_' + targetLang;

                $('#languagePicker li').removeClass('active');
                $(this).parent().addClass('active');

                selectedLanguage = targetLang;

                $(targetDivID).toggleClass('in').siblings().removeClass('in');

                setLanguagePicker(selectedLanguage);
                setNavbarLang(selectedLanguage);
            }
        );


    });
}());