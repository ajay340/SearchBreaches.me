const search = document.getElementById('searchBox');
const matchList = document.getElementById('searchResults');
let breaches;


   const getBreachSearch = async (searchText) =>{
    let matches = [];
    if (searchText.length === 0) {
        matches = [];
        matchList.innerHTML = '';
        matchList.className = "term-list hidden";
    } else{
    url = '/searchq=' + searchText
    const res = await fetch(url);
    breaches = await res.json();
    matches = breaches.filter(breaches => {
        const regex = new RegExp(`^${searchText}`, 'gi');
        return breaches.Name_of_Covered_Entity.match(regex) || breaches.Industry.match(regex) || breaches.Type_of_Breach.match(regex); 
       });
      
       // Clear when input or matches are empty
       if (searchText.length === 0) {
        matches = [];
        matchList.innerHTML = '';
        matchList.className = "term-list hidden";
       }
    }
      
       outputHtml(matches);
      };

const outputHtml = matches => {
 if (matches.length > 0) {
  const html = matches
   .map(
    match => `<li><a href="/product/${match.ID}">${match.Name_of_Covered_Entity}</a></li>`
   )
   .join('');
  matchList.innerHTML = html;
  matchList.className = "term-list"
 }
};

search.addEventListener('input', () => getBreachSearch(search.value));