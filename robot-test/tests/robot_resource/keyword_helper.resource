*** Settings ***

Library    RequestsLibrary


*** Variables ***

${EVENT_PRODUCER_HOST}         event-producer-svc
${EVENT_PRODUCER_PORT}         8080
${EVENT_PRODUCER_URL}          http://${EVENT_PRODUCER_HOST}:${EVENT_PRODUCER_PORT}
${TEST_NAMESPACE}              default
${REDIS_POD_LABEL}             app.kubernetes.io/name=redis,app.kubernetes.io/component=master
${EVENT_PRODUCER_POD_LABEL}    app=event-producer
${REDIS_SECRET_NAME}           redis
${REDIS_PORT}                  6379
${REDIS_STREAM_NAME}           NotificationStream


*** Keywords ***

Send Event Message
    [Documentation]    Event Notification (EN) interface
    [Arguments]
    ...    ${data}
    ...    ${expected_status}

    POST
    ...    url=${EVENT_PRODUCER_URL}/event
    ...    data=${data}
    ...    expected_status=${expected_status}

Send Invalid Event Message
        [Documentation]    Event Notification (EN) interface invalid request
        [Arguments]
        ...    ${data}
        ...    ${expected_status}
        ...    ${is_post}=${True}

    IF    ${is_post}
        POST
        ...    url=${EVENT_PRODUCER_URL}/event
        ...    data=${data}
        ...    expected_status=${expected_status}
    ELSE
        PUT
        ...    url=${EVENT_PRODUCER_URL}/event
        ...    data=${data}
        ...    expected_status=${expected_status}
    END