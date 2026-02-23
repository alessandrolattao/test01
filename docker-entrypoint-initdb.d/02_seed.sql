-- Seed data for development

-- Admin user: admin@test.com / password
-- bcrypt hash of "password" with cost 10
INSERT INTO admins (id, email, password_hash) VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'admin@test.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'
);

-- Psychological Profile Questionnaire (Big Five / OCEAN)
-- 10 situational questions, 2 per dimension, max score 100
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
    -- OPENNESS (questions 1-2)
    -- =============================================

    -- Q1: Openness - Reaction to new technologies
    q1 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q1, qid,
        'Il tuo team deve affrontare un progetto completamente nuovo, con tecnologie e metodi che nessuno ha mai usato prima. Come reagisci?',
        1
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q1, 'Sono entusiasta: mi offro volontario per esplorare le nuove tecnologie e propongo un piano di studio condiviso', 10, 1),
        (q1, 'Sono curioso ma cauto: inizio a documentarmi per conto mio prima di espormi con il team', 7, 2),
        (q1, 'Preferisco che qualcun altro faccia da apripista, poi mi adeguo una volta che il percorso e'' chiaro', 4, 3),
        (q1, 'Propongo di usare metodi gia'' collaudati e adattarli al nuovo progetto, evitando rischi inutili', 2, 4);

    -- Q2: Openness - Handling opposing ideas
    q2 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q2, qid,
        'Durante una riunione, un collega presenta un''idea che contraddice completamente la tua visione del progetto. Come ti comporti?',
        2
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q2, 'Ascolto con attenzione e chiedo di approfondire: potrebbe avere visto qualcosa che mi e'' sfuggito', 10, 1),
        (q2, 'Prendo nota della sua idea e la confronto con la mia dopo la riunione, con calma', 7, 2),
        (q2, 'Espongo subito i punti deboli della sua proposta per evitare che il team prenda una direzione sbagliata', 3, 3),
        (q2, 'Lascio correre per non creare conflitto, ma resto convinto della mia idea originale', 1, 4);

    -- =============================================
    -- CONSCIENTIOUSNESS (questions 3-4)
    -- =============================================

    -- Q3: Conscientiousness - Deadline management
    q3 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q3, qid,
        'Hai tre task da completare entro venerdi''. Il martedi'' ti accorgi che uno dei tre richiede piu'' tempo del previsto. Cosa fai?',
        3
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q3, 'Ricalcolo le priorita'', comunico il ritardo al responsabile e propongo un piano aggiornato con le nuove tempistiche', 10, 1),
        (q3, 'Mi organizzo per fare straordinario e cerco di consegnare tutto nei tempi, anche a costo di sacrificare qualcosa nella qualita''', 6, 2),
        (q3, 'Mi concentro sul task problematico e spero di riuscire a recuperare sugli altri nei giorni restanti', 4, 3),
        (q3, 'Aspetto di vedere come evolve la situazione: magari il task si risolve piu'' in fretta di quanto penso', 1, 4);

    -- Q4: Conscientiousness - Error handling integrity
    q4 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q4, qid,
        'Stai lavorando a un report importante. A meta'' del lavoro ti rendi conto che hai commesso un errore nei dati iniziali che nessuno noterebbe. Come ti comporti?',
        4
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q4, 'Mi fermo, correggo l''errore dall''inizio e ricalcolo tutto, anche se significa perdere mezza giornata di lavoro', 10, 1),
        (q4, 'Correggo i dati da quel punto in poi e aggiungo una nota interna che spiega la correzione', 7, 2),
        (q4, 'Valuto l''impatto dell''errore: se e'' minimo, vado avanti e lo segnalo dopo la consegna', 4, 3),
        (q4, 'Se nessuno lo noterebbe, non vale la pena perdere tempo: consegno nei tempi e passo oltre', 1, 4);

    -- =============================================
    -- EXTRAVERSION (questions 5-6)
    -- =============================================

    -- Q5: Extraversion - Networking event
    q5 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q5, qid,
        'L''azienda organizza un evento di networking con professionisti del tuo settore che non conosci. Come affronti la serata?',
        5
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q5, 'Mi presento a piu'' persone possibile, cerco punti in comune e scambio contatti durante tutta la serata', 10, 1),
        (q5, 'Mi inserisco nei gruppi di conversazione gia'' formati e contribuisco quando ho qualcosa di interessante da dire', 7, 2),
        (q5, 'Cerco qualcuno che conosco gia'' e resto nella sua orbita, aprendomi gradualmente ad altri', 4, 3),
        (q5, 'Partecipo ma resto in disparte, osservo e intervengo solo se qualcuno mi coinvolge direttamente', 2, 4);

    -- Q6: Extraversion - Team strategy discussion
    q6 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q6, qid,
        'Il tuo team deve decidere la strategia per il prossimo trimestre. Il responsabile chiede opinioni a tutti. Qual e'' il tuo approccio?',
        6
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q6, 'Prendo la parola tra i primi, espongo la mia visione e cerco di guidare la discussione verso una decisione condivisa', 10, 1),
        (q6, 'Ascolto le prime proposte, poi intervengo con la mia opinione arricchita dai contributi degli altri', 8, 2),
        (q6, 'Preparo i miei punti ma aspetto che il responsabile mi chieda direttamente cosa ne penso', 4, 3),
        (q6, 'Preferisco mandare le mie idee via email dopo la riunione, quando ho avuto tempo di elaborarle meglio', 2, 4);

    -- =============================================
    -- AGREEABLENESS (questions 7-8)
    -- =============================================

    -- Q7: Agreeableness - Struggling colleague
    q7 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q7, qid,
        'Un collega con cui lavori a stretto contatto sta attraversando un periodo difficile e le sue prestazioni ne risentono, aumentando il tuo carico di lavoro. Come gestisci la situazione?',
        7
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q7, 'Mi offro di aiutarlo con i suoi task e gli chiedo se c''e'' qualcosa che posso fare per alleggerirgli il periodo', 10, 1),
        (q7, 'Assorbo il carico extra senza dire nulla, ma cerco di capire discretamente come sta', 7, 2),
        (q7, 'Gestisco il mio lavoro normalmente e segnalo al responsabile che c''e'' un problema di capacita'' nel team', 4, 3),
        (q7, 'Gli faccio presente che il suo calo sta impattando il mio lavoro e gli chiedo di trovare una soluzione', 2, 4);

    -- Q8: Agreeableness - Team conflict mediation
    q8 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q8, qid,
        'Durante un progetto di gruppo, emerge un conflitto tra due colleghi che blocca l''avanzamento del lavoro. Tu non sei direttamente coinvolto. Cosa fai?',
        8
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q8, 'Mi propongo come mediatore: parlo con entrambi separatamente per capire le posizioni e cerco un compromesso', 10, 1),
        (q8, 'Suggerisco al team di fare una riunione dove ognuno espone il suo punto di vista, per trovare una soluzione insieme', 8, 2),
        (q8, 'Lascio che risolvano tra loro: non e'' il mio conflitto e non voglio essere coinvolto', 3, 3),
        (q8, 'Segnalo la situazione al responsabile perche'' intervenga prima che il progetto ne risenta troppo', 5, 4);

    -- =============================================
    -- EMOTIONAL STABILITY (questions 9-10)
    -- =============================================

    -- Q9: Emotional Stability - Last-minute changes
    q9 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q9, qid,
        'Il giorno prima di una presentazione importante, scopri che il cliente ha cambiato completamente le richieste. Il tuo lavoro degli ultimi due giorni e'' da rifare. Qual e'' la tua prima reazione?',
        9
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q9, 'Faccio un respiro, valuto cosa posso recuperare dal lavoro gia'' fatto e mi organizzo per adattare la presentazione', 10, 1),
        (q9, 'Sono frustrato ma mi metto subito al lavoro: lamentarmi non cambierebbe la situazione', 7, 2),
        (q9, 'Mi sfogo con un collega per scaricare la tensione, poi mi rimetto al lavoro con piu'' calma', 5, 3),
        (q9, 'Sento la pressione salire e faccio fatica a concentrarmi: il cambio dell''ultimo minuto mi manda in confusione', 2, 4);

    -- Q10: Emotional Stability - Negative feedback
    q10 := gen_random_uuid();
    INSERT INTO questions (id, questionnaire_id, text, sort_order) VALUES (
        q10, qid,
        'Il tuo responsabile ti da'' un feedback negativo su un progetto in cui avevi investito molto impegno. Come reagisci nelle ore successive?',
        10
    );
    INSERT INTO answers (question_id, text, score, sort_order) VALUES
        (q10, 'Analizzo il feedback punto per punto, identifico cosa migliorare e preparo un piano d''azione per il prossimo progetto', 10, 1),
        (q10, 'Mi prendo un momento per digerire la delusione, poi rileggo il feedback con occhi piu'' oggettivi', 7, 2),
        (q10, 'Ci rimango male e ci penso per il resto della giornata, ma so che passera''', 4, 3),
        (q10, 'Metto in discussione le mie capacita'': forse questo lavoro non fa per me', 1, 4);

END $$;
