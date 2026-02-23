-- Seed data for development

-- Admin user: admin@test.com / password
-- bcrypt hash of "password" with cost 10
INSERT INTO admins (id, email, password_hash) VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'admin@test.com',
    '$2b$10$suJ3c8ZWOJq/xwcTWa7QquJGSgAnE54D8pNQx/if60kjDBKo9cZPi'
);

-- Tourism Reviewer Selector Questionnaire v1
-- 10 situational questions, 2 per dimension, max score 100
-- Dimensions: Attenzione al Dettaglio, Comunicazione Efficace,
--             Obiettivita e Imparzialita, Autonomia e Organizzazione,
--             Sensibilita Culturale
INSERT INTO questionnaires (id, version, is_active) VALUES (
    gen_random_uuid(), 1, TRUE
);

-- Store questionnaire ID for FK references
DO $$
DECLARE
    qid UUID;
    q1 UUID; q2 UUID; q3 UUID; q4 UUID; q5 UUID;
    q6 UUID; q7 UUID; q8 UUID; q9 UUID; q10 UUID;
BEGIN
    SELECT id INTO qid FROM questionnaires WHERE is_active = TRUE LIMIT 1;

    -- =============================================
    -- ATTENZIONE AL DETTAGLIO (questions 1-2)
    -- =============================================

    -- Q1: Attention to Detail - Restaurant review completeness
    q1 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q1, qid,
        'Stai visitando un ristorante per recensirlo. Il cibo e'' buono ma noti che il bagno non e'' pulitissimo e l''illuminazione del locale e'' troppo forte. Nella tua recensione, come tratti questi aspetti?',
        1
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q1, 'Li descrivo entrambi con precisione, spiegando come impattano l''esperienza complessiva del cliente, pur riconoscendo la qualita'' del cibo', 10, 1),
        (q1, 'Menziono il bagno perche'' riguarda l''igiene, ma tralascio l''illuminazione che e'' una questione di gusto personale', 6, 2),
        (q1, 'Mi concentro sul cibo che e'' l''aspetto principale: i dettagli secondari li segnalo solo se sono davvero gravi', 3, 3),
        (q1, 'Non li riporto: se il cibo e'' buono, il resto e'' contorno e non voglio penalizzare il ristorante per dettagli minori', 1, 4);

    -- Q2: Attention to Detail - Unlisted museum discovery
    q2 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q2, qid,
        'Stai censendo le attivita'' di un borgo turistico e noti che un piccolo museo non e'' presente su nessuna guida online. All''ingresso non ci sono orari visibili e il custode non e'' disponibile. Come procedi?',
        2
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q2, 'Fotografo l''esterno, annoto l''indirizzo esatto, cerco un contatto telefonico sui cartelli, chiedo informazioni ai commercianti vicini e torno in un altro momento per completare la scheda', 10, 1),
        (q2, 'Segno il nome e l''indirizzo, scatto una foto e lo inserisco nella lista come attivita'' da verificare in seguito', 7, 2),
        (q2, 'Lo inserisco nella lista con le informazioni che riesco a trovare online, senza tornare di persona', 4, 3),
        (q2, 'Se non c''e'' online e non riesco a entrare, probabilmente non e'' un''attivita'' rilevante per i turisti: passo oltre', 1, 4);

    -- =============================================
    -- COMUNICAZIONE EFFICACE (questions 3-4)
    -- =============================================

    -- Q3: Effective Communication - Wine tasting review
    q3 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q3, qid,
        'Devi descrivere un''esperienza di degustazione di vini locali. L''evento era piacevole ma non eccezionale. Come imposti la recensione?',
        3
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q3, 'Racconto l''esperienza in modo coinvolgente partendo dall''atmosfera, descrivo i vini con dettagli sensoriali accessibili anche ai non esperti, e chiudo con un giudizio onesto su per chi e'' consigliata', 10, 1),
        (q3, 'Descrivo i vini provati, il prezzo e la durata dell''esperienza, aggiungendo le mie impressioni personali in modo chiaro e diretto', 7, 2),
        (q3, 'Scrivo una recensione tecnica dettagliata sui vini, con note sui vitigni e le tecniche di produzione, per dare informazioni precise', 4, 3),
        (q3, 'Scrivo che e'' stata un''esperienza carina ma niente di speciale, con un voto medio e qualche riga di commento', 2, 4);

    -- Q4: Effective Communication - Luxury hotel review
    q4 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q4, qid,
        'Ti chiedono di recensire un hotel di lusso. Non hai mai soggiornato in strutture di questa categoria. Come affronti la scrittura?',
        4
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q4, 'Studio le recensioni di riferimento per capire i parametri del segmento lusso, vivo l''esperienza con attenzione a ogni aspetto e scrivo la recensione bilanciando descrizione oggettiva e impressioni da ospite, specificando il mio punto di vista', 10, 1),
        (q4, 'Mi informo sui servizi standard delle strutture di lusso per avere un metro di paragone e scrivo la recensione confrontando quello che trovo con quello che ci si aspetta', 7, 2),
        (q4, 'Scrivo quello che vedo e che provo: la mia inesperienza nel lusso puo'' essere un vantaggio perche'' rappresento il turista medio', 5, 3),
        (q4, 'Chiedo a un collega piu'' esperto di strutture di lusso di aiutarmi con la recensione o di farsene carico', 2, 4);

    -- =============================================
    -- OBIETTIVITA E IMPARZIALITA (questions 5-6)
    -- =============================================

    -- Q5: Objectivity - Free meal offer
    q5 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q5, qid,
        'Il proprietario di un agriturismo che stai recensendo ti offre un pranzo gratuito e una bottiglia di vino locale in regalo. Come ti comporti?',
        5
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q5, 'Ringrazio cortesemente ma declino sia il pranzo che il regalo, spiegando che devo mantenere l''imparzialita'' della mia valutazione. Pago il mio pranzo regolarmente', 10, 1),
        (q5, 'Accetto il pranzo perche'' fa parte dell''esperienza da valutare, ma rifiuto il vino in regalo. Nella recensione specifico che sono stato ospite', 7, 2),
        (q5, 'Accetto con piacere entrambi: e'' un gesto di ospitalita'' tipico della cultura locale e sarebbe scortese rifiutare. Cerco comunque di essere obiettivo nella recensione', 3, 3),
        (q5, 'Accetto tutto senza problemi: e'' normale che le attivita'' cerchino di fare bella figura con i recensori, tanto so essere obiettivo a prescindere', 1, 4);

    -- Q6: Objectivity - Bad experience vs positive reviews
    q6 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q6, qid,
        'Stai recensendo un ristorante e la tua esperienza e'' stata pessima: servizio lento, piatti tiepidi, conto sbagliato. Pero'' online ha centinaia di recensioni positive. Come scrivi la tua valutazione?',
        6
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q6, 'Descrivo la mia esperienza con fatti specifici (tempi di attesa, temperatura dei piatti, errore nel conto), riconosco che le recensioni generali sono positive e suggerisco che potrebbe trattarsi di una giornata sfortunata, senza pero'' alterare il mio giudizio', 10, 1),
        (q6, 'Scrivo esattamente quello che ho vissuto senza farmi influenzare dalle altre recensioni: la mia esperienza e'' quella che conta per la mia valutazione', 7, 2),
        (q6, 'Ammorbidisco un po'' il giudizio considerando le tante recensioni positive: forse ho beccato un giorno storto e non sarebbe giusto penalizzarlo troppo', 3, 3),
        (q6, 'Con cosi'' tante recensioni positive, forse il problema sono le mie aspettative. Riscrivo la recensione dando piu'' peso agli aspetti positivi', 1, 4);

    -- =============================================
    -- AUTONOMIA E ORGANIZZAZIONE (questions 7-8)
    -- =============================================

    -- Q7: Autonomy - Week planning for 15 activities
    q7 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q7, qid,
        'Hai una settimana per censire e recensire 15 attivita'' sparse in una zona collinare con trasporti pubblici limitati. Come organizzi il lavoro?',
        7
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q7, 'Mappo tutte le attivita'', le raggruppo per zona geografica, pianifico un percorso ottimale giorno per giorno verificando orari di apertura e accessibilita'', e tengo due slot liberi per imprevisti o attivita'' scoperte sul posto', 10, 1),
        (q7, 'Organizzo il percorso per aree geografiche in modo da minimizzare gli spostamenti, verifico gli orari di apertura e prenoto dove necessario', 7, 2),
        (q7, 'Faccio un piano generale dividendo le 15 attivita'' nei 7 giorni, ma preferisco adattarmi giorno per giorno in base a come procede il lavoro', 4, 3),
        (q7, 'Inizio dalla prima attivita'' della lista e procedo in ordine: pianificare troppo in anticipo e'' inutile perche'' sul campo cambia sempre tutto', 1, 4);

    -- Q8: Autonomy - Falling behind schedule
    q8 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q8, qid,
        'E'' il terzo giorno di lavoro sul territorio e ti rendi conto che il ritmo che hai tenuto finora non ti permettera'' di completare tutte le attivita'' assegnate entro la scadenza. Cosa fai?',
        8
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q8, 'Ricalcolo il piano: identifico le attivita'' prioritarie, comunico subito al coordinatore il potenziale ritardo e propongo un piano B con una nuova timeline per quelle restanti', 10, 1),
        (q8, 'Aumento il ritmo, accetto di fare recensioni piu'' sintetiche per le attivita'' meno importanti e cerco di portare a termine tutto nei tempi', 5, 2),
        (q8, 'Comunico al coordinatore che non riusciro'' a finire tutto e chiedo come procedere, aspettando istruzioni', 4, 3),
        (q8, 'Continuo con il mio ritmo e faccio quello che riesco: meglio poche recensioni fatte bene che tante fatte di corsa', 3, 4);

    -- =============================================
    -- SENSIBILITA CULTURALE (questions 9-10)
    -- =============================================

    -- Q9: Cultural Sensitivity - Temple visit
    q9 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q9, qid,
        'Stai visitando un tempio di una religione che non conosci bene. I visitatori locali seguono rituali specifici all''ingresso. Non sai se siano obbligatori anche per i turisti. Come ti comporti?',
        9
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q9, 'Osservo con attenzione cosa fanno gli altri visitatori, chiedo a qualcuno del posto se ci sono regole da seguire e mi adeguo rispettosamente. Nella recensione includo indicazioni pratiche per i turisti sui comportamenti attesi', 10, 1),
        (q9, 'Seguo quello che fanno gli altri per rispetto, anche senza capirne il significato. Nella recensione menziono i rituali e consiglio ai turisti di informarsi prima della visita', 7, 2),
        (q9, 'Entro rispettosamente senza partecipare ai rituali: non essendo della loro religione, non voglio essere inappropriato nel praticare riti che non conosco', 4, 3),
        (q9, 'Entro come visitatore normale: sono un turista, non un fedele. I rituali religiosi sono per chi pratica quella fede', 1, 4);

    -- Q10: Cultural Sensitivity - Adapted ethnic restaurant
    q10 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q10, qid,
        'Stai recensendo un ristorante etnico in una zona turistica. I piatti sono stati adattati al gusto locale e sono diversi dalla versione originale della cucina. Come valuti questo aspetto nella recensione?',
        10
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q10, 'Descrivo i piatti per quello che sono, segnalo le differenze rispetto alla tradizione originale in modo informativo (non giudicante), e valuto la qualita'' dell''esecuzione indipendentemente dall''autenticita''', 10, 1),
        (q10, 'Valuto i piatti in base alla qualita'' complessiva e al rapporto qualita''-prezzo. Menziono che la cucina e'' adattata al gusto locale, lasciando al lettore il giudizio', 7, 2),
        (q10, 'Segnalo chiaramente che non si tratta di cucina autentica: i turisti devono sapere che i piatti sono stati modificati rispetto all''originale', 4, 3),
        (q10, 'Penalizzo il ristorante nella valutazione: se proponi una cucina etnica, deve essere autentica, altrimenti e'' una mancanza di rispetto per la cultura originale', 1, 4);

END $$;
