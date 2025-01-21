package flags

import "regexp"

var flaggedPhrases = []*regexp.Regexp{
	regexp.MustCompile(`(?mi).*REPOST THIS VIDEO AND I WILL REPOST 2 YOURS.*`),
	regexp.MustCompile(`(?mis).*All i did was simple online work from(?s).*`),
	regexp.MustCompile(`(?mis).*hr provide by Google.*`),
	regexp.MustCompile(`(?mis).*hour for doing online work from home(?s).*`),
	regexp.MustCompile(`(?mis).*adshrink.it.*`),
	regexp.MustCompile(`(?mis).*I get paid over \$.* per hour working from home(?s).*`),
	regexp.MustCompile(`(?mis).*I get paid more than \$(?s).* per hour working from home(?s).*`),
	regexp.MustCompile(`(?mis).*I get paid more than \$(?s).* per hour working online(?s).*`),
	regexp.MustCompile(`(?mis).*Real online home based job(?s).*`),
	regexp.MustCompile(`(?mis).*I getting Paid upto \$.* this week(?s).*`),
	regexp.MustCompile(`(?mis).*Get your financial freedom out of the hands(?s).*`),
	regexp.MustCompile(`(?mis).*simple online job from home(?s).*`),
	regexp.MustCompile(`(?mis).*She has been fired for .* months but last month(?s).*`),
	regexp.MustCompile(`(?mis).*She has been laid off for .* months but last(?s).*`),
	regexp.MustCompile(`(?mis).*makes \$.*hour.* on the computer(?s).*`),
	regexp.MustCompile(`(?mis).*makes \$.*hour.* on the laptop(?s).*`),
	regexp.MustCompile(`(?mis).*makes \$.*hour.* on the internet(?s).*`),
	regexp.MustCompile(`(?mis).*makes \$.*hour.* from home(?s).*`),
	regexp.MustCompile(`(?mis).*peau de mouton qui habille et camoufle un loup.*`),
	regexp.MustCompile(`(?mis).*\$.*-\$.* per month on the web(?s).*`),
	regexp.MustCompile(`(?mi).*Help the channel and join you these awesome pages to obtain free crypto coins.*`),
	regexp.MustCompile(`(?mis).*per hour for doing work online work(?s).*`),
	regexp.MustCompile(`(?mis).*earned \$.* in my first .* month(?s).*`),
	regexp.MustCompile(`(?mi).*Porn1MinHD-LBRY.*`),
	regexp.MustCompile(`(?mis).*Do you want to make money on your phone?(?s).*`),
	regexp.MustCompile(`(?mi).*sigo a todos los que.*`),
	regexp.MustCompile(`(?mi).*Comentá para seguir comentando tu publicación.*`),
	regexp.MustCompile(`(?mi).*BRASILemEVID.*`),
	regexp.MustCompile(`(?mi).*Sígueme de vuelta.*`),
	regexp.MustCompile(`(?mi).*Necessito de seguidores.*`),
	regexp.MustCompile(`(?mi).*sigo y me siguen.*`),
	regexp.MustCompile(`(?mi).*follow me as I follow you.*`),
	regexp.MustCompile(`(?mi).*@PremiumPorn.*`),
	regexp.MustCompile(`(?mi).*@EvaElfie.*`),
	regexp.MustCompile(`(?mi).*@niggertown.*`),
	regexp.MustCompile(`(?mis).*generating extra cash online from home more(?s).*`),
	regexp.MustCompile(`(?mis).*for more info visit any tab this site(?s).*`),
	regexp.MustCompile(`(?mis).*started earning \$.* hour.* in my free time(?s).*`),
	regexp.MustCompile(`(?mis).*making over \$.* a month working part time(?s).*`),
	regexp.MustCompile(`(?mis).*making about \$.*-\$.* per month and you can too(?s).*`),
	regexp.MustCompile(`(?mis).*paid over(?s).*https://www.*`),
	regexp.MustCompile(`(?mis).*working from(?s).*https://www.*`),
	regexp.MustCompile(`(?mis).*working at home(?s).*https://www.*`),
	regexp.MustCompile(`(?mis).*per-hr(?s).*earning $(?s).*http*`),
	regexp.MustCompile(`(?mis).*paid more(?s).*https://tinyurl*`),
	regexp.MustCompile(`(?mis).*earn over(?s).*https://tinyurl*`),
	regexp.MustCompile(`(?mis).*I recommend him to anyone with; Herpes Virus.*`),
	regexp.MustCompile(`(?mis).*Nicholas Jonathan Gregory.*`),
	regexp.MustCompile(`(?mis).*Windsor Dr..*`),
	regexp.MustCompile(`(?mis).*SunlightAFA.*`),
}
