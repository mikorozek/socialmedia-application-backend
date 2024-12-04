pipeline {
    agent any

    environment {
        NEXUS_URL = 'http://52.233.173.205:8081'  // Adres Twojego Nexusa
        NEXUS_REPO = 'maven-releases'  // Repozytorium w Nexusie, gdzie będziesz wysyłać artefakty
        NEXUS_CREDENTIALS = 'nexus-admin'  // ID poświadczeń w Jenkinsie
        MAVEN_HOME = '/usr/share/maven'   // Ścieżka do Mavena w kontenerze Docker (jeśli potrzebne)
    }

    stages {
        // Etap pobierania kodu z repozytorium
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        // Etap budowania aplikacji Spring Boot przy użyciu Dockera
        stage('Build') {
            agent {
                docker {
                    image 'maven:3.8.4-openjdk-11'  // Użycie obrazu Dockera z Maven i OpenJDK 11
                    args '-v $HOME/.m2:/root/.m2'  // Montowanie katalogu z lokalnymi repozytoriami Maven
                }
            }
            steps {
                script {
                    // Używamy Maven do budowy aplikacji
                    sh 'mvn clean package'
                }
            }
        }

        // Etap testowania aplikacji (opcjonalny)
        stage('Test') {
            agent {
                docker {
                    image 'maven:3.8.4-openjdk-11'  // Obraz z Maven i OpenJDK 11 do uruchamiania testów
                    args '-v $HOME/.m2:/root/.m2'  // Montowanie katalogu z lokalnymi repozytoriami Maven
                }
            }
            steps {
                script {
                    // Uruchamianie testów jednostkowych
                    sh 'mvn test'
                }
            }
        }

        // Etap przesyłania artefaktów do Nexusa
        stage('Upload to Nexus') {
            agent {
                docker {
                    image 'maven:3.8.4-openjdk-11'  // Obraz z Maven do przesyłania artefaktów
                    args '-v $HOME/.m2:/root/.m2'  // Montowanie katalogu z lokalnymi repozytoriami Maven
                }
            }
            steps {
                script {
                    // Używamy Mavena do przesyłania artefaktów do Nexusa
                    sh "mvn deploy -DskipTests -DaltDeploymentRepository=nexus::default::${NEXUS_URL}/repository/${NEXUS_REPO}"
                }
            }
        }
    }
}
