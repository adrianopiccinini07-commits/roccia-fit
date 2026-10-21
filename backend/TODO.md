# Roadmap di Sviluppo: Roccia Fit

## Fase 1: Autenticazione e Sicurezza Client (Frontend)
- [ ] Installare `expo-secure-store` nel progetto React Native.
- [ ] Aggiornare `login.tsx` per salvare il JWT in modo crittografato nel dispositivo.
- [ ] Implementare un "Auth Context" (Context API o Zustand) per gestire lo stato globale (utente loggato/non loggato).
- [ ] Configurare il routing protetto: se c'è il token vai alla Dashboard, altrimenti vai al Login.
- [ ] Creare un pulsante e una funzione di "Logout" (che cancelli il token e riavvii il ciclo).

## Fase 2: Blindare il Server (Backend Go)
- [ ] Creare un pacchetto `middleware` nel backend Go.
- [ ] Scrivere un middleware di autenticazione che intercetti l'header `Authorization: Bearer <token>`.
- [ ] Scrivere la logica per decodificare e validare la firma e la scadenza del JWT.
- [ ] Applicare questo middleware nel router `chi` a un nuovo gruppo di rotte (es. `/api/v1/private/...`).
- [ ] Creare l'endpoint `GET /api/v1/users/me` per verificare che l'utente loggato sia riconosciuto dal server.

## Fase 3: Funzionalità Core "Roccia Fit" (Backend Go)
- [ ] Progettare le tabelle GORM per i dati fisici (es. log del peso giornaliero, circonferenze, pliche per preparazioni Classic Physique).
- [ ] Scrivere i controller per salvare un nuovo log nel database (Create).
- [ ] Scrivere i controller per recuperare lo storico dei dati di un utente (Read).
- [ ] Testare accuratamente i nuovi endpoint usando chiamate API autorizzate con il JWT.

## Fase 4: Interfaccia e Dashboard (Frontend React Native)
- [ ] Creare la pagina `dashboard.tsx` che farà da Home dell'app.
- [ ] Programmare la chiamata HTTP (con fetch) che invii il token al backend per richiedere i dati utente.
- [ ] Disegnare un modulo/form per l'inserimento rapido dei parametri (peso, note sull'allenamento).
- [ ] Integrare una libreria di grafici (es. `react-native-chart-kit` o `react-native-gifted-charts`) per tracciare visivamente i progressi nel tempo.

## Fase 5: Preparazione alla Produzione (Deploy)
- [ ] Configurare le variabili d'ambiente (`.env`) su frontend e backend per nascondere IP, chiavi segrete e password del DB.
- [ ] Migrare da Docker locale a un database PostgreSQL in cloud (es. Supabase, Render o AWS).
- [ ] Configurare l'hosting per il server Go in modo che sia sempre raggiungibile su internet (VPS o PaaS).
- [ ] Sostituire `10.0.2.2` nell'app React Native con l'URL di produzione reale.
- [ ] Generare la prima build nativa Android/iOS tramite Expo EAS (`eas build`) per installarla sul telefono fisico.