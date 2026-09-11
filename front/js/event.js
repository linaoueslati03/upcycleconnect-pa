document.querySelectorAll('.pill').forEach(function(bouton) {

    bouton.addEventListener('click', function() {

        const statut = bouton.dataset.statut;

        document.querySelectorAll('.pill').forEach(function(b) {
            b.classList.remove('active');
        });

        bouton.classList.add('active');

        document.querySelectorAll('.evenement-row').forEach(function(ligne) {

            if (statut === 'tous' || ligne.dataset.statut === statut) {
                ligne.style.display = '';
            } else {
                ligne.style.display = 'none';
            }

        });

    });

});