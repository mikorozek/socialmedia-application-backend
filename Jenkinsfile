pipeline {
    agent any

    environment {
        NEXUS_URL = 'http://52.233.173.205:8081'
        NEXUS_REPO = 'maven-repository'
        NEXUS_CREDENTIALS = 'nexus-admin2'
        MAVEN_HOME = '/usr/share/maven'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build') {
            agent {
                docker {
                    image 'maven:3.8.4-openjdk-11'
                    args '-v $HOME/.m2:/root/.m2'
                }
            }
            steps {
                script {
                    sh 'mvn clean package'
                }
            }
        }

        stage('Test') {
            agent {
                docker {
                    image 'maven:3.8.4-openjdk-11'
                    args '-v $HOME/.m2:/root/.m2'
                }
            }
            steps {
                script {
                    sh 'mvn test'
                }
            }
        }

        stage('Upload to Nexus') {
            agent {
                docker {
                    image 'maven:3.8.4-openjdk-11'
                    args '-v $HOME/.m2:/root/.m2'
                }
            }
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: 'nexus-admin', usernameVariable: 'NEXUS_USERNAME', passwordVariable: 'NEXUS_PASSWORD')]) {
                        writeFile file: 'settings.xml', text: """
                            <settings>
                                <servers>
                                    <server>
                                        <id>nexus</id>
                                        <username>${NEXUS_USERNAME}</username>
                                        <password>${NEXUS_PASSWORD}</password>
                                    </server>
                                </servers>
                            </settings>
                        """
                        sh "mvn clean deploy -s settings.xml -DskipTests -DaltDeploymentRepository=nexus::default::${NEXUS_URL}/repository/${NEXUS_REPO}"
                    }
                }
            }
        }
    }
}
