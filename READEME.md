# SwarmGuardian

Ao trabalhar em uma equipe que administrava um cluster Docker Swarm, percebemos que delegar o controle de stacks para os desenvolvedores, as vezes, poderia ser um pouco perigoso. Em um caso isolado um analista acabou subindo um service com 30 replicas durante o horario de pico de acessos de nossos sistemas!

Rapidamente percebemos que o Docker Swarm não tinha nenhuma forma de evitar este tipo de caso, sem retirar completamente o acesso de terceiros a nosso ambiente.

Com isto criamos o SwarmGuardian. Usando as informações coletadas pelo docker-flow-swarm-listener (https://github.com/vfarcic/docker-flow-swarm-listener), ele verificará se os serviços criados/atualizados tem mais replicas do que um valor definido. Caso isto aconteça, nos enviará uma notificação no Slack.

Ele foi desenvolvido GoLang. Ele é bem simples, porém com o tempo esperamos expandi-lo.

## Modo de uso

Este docker-compose pode ser usado para publicar o SwarmGuardian.

version: "3.3"

services:

  swarm-listener:
    image: dockerflow/docker-flow-swarm-listener:18.11.28-19
    environment:
      DF_NOTIFY_CREATE_SERVICE_URL: http://notifier:8081/reconfigure
      DF_NOTIFY_LABEL: com.docker.stack.namespace
      DF_SERVICE_POLLING_INTERVAL: 20
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    deploy:
      mode: replicated
      placement:
        constraints:
          - "node.role==manager"
      restart_policy:
        condition: any
    networks:
      - swarmguardian

  notifier:
    image: jarzamendia/swarmguardian:v1
    networks:
      - swarmguardian

networks:
  swarmguardian: