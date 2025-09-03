CREATE TABLE IF NOT EXISTS contacorrente (
	idcontacorrente TEXT(37) PRIMARY KEY,
	numero INTEGER(10) NOT NULL UNIQUE,
	nome TEXT(100) NOT NULL,
	cpf TEXT(11) NOT NULL UNIQUE,
	ativo INTEGER(1) NOT NULL default 1,
	senha TEXT(100) NOT NULL,
	salt TEXT(100) NOT NULL,
	CHECK (ativo in (0,1))
);

CREATE TABLE IF NOT EXISTS movimento (
	idmovimento TEXT(37) PRIMARY KEY,
	idcontacorrente TEXT(37) NOT NULL,
	datamovimento TEXT(25) NOT NULL,
	tipomovimento TEXT(1) NOT NULL,
	valor REAL NOT NULL,
	idempotencia_key TEXT(37),
	CHECK (tipomovimento in ('C','D')),
	FOREIGN KEY(idcontacorrente) REFERENCES contacorrente(idcontacorrente)
);

CREATE TABLE IF NOT EXISTS tarifa (
	idtarifa TEXT(37) PRIMARY KEY,
	idcontacorrente TEXT(37) NOT NULL,
	datamovimento TEXT(25) NOT NULL,
	valor REAL NOT NULL,
	FOREIGN KEY(idcontacorrente) REFERENCES contacorrente(idcontacorrente)
);

CREATE TABLE IF NOT EXISTS transferencia (
	idtransferencia TEXT(37) PRIMARY KEY,
	idcontacorrente_origem TEXT(37) NOT NULL,
	idcontacorrente_destino TEXT(37) NOT NULL,
	datamovimento TEXT(25) NOT NULL,
	valor REAL NOT NULL,
	status INTEGER(1) NOT NULL DEFAULT 0,
	data_conclusao TEXT(25),
	descricao TEXT(255),
	idempotencia_key TEXT(37),
	CHECK (status in (0,1,2)),
	FOREIGN KEY(idcontacorrente_origem) REFERENCES contacorrente(idcontacorrente),
	FOREIGN KEY(idcontacorrente_destino) REFERENCES contacorrente(idcontacorrente)
);

CREATE TABLE IF NOT EXISTS idempotencia (
	chave_idempotencia TEXT(37) PRIMARY KEY,
	requisicao TEXT(1000),
	resultado TEXT(1000)
);

CREATE INDEX IF NOT EXISTS idx_movimento_conta ON movimento(idcontacorrente);
CREATE INDEX IF NOT EXISTS idx_transferencia_origem ON transferencia(idcontacorrente_origem);
CREATE INDEX IF NOT EXISTS idx_transferencia_destino ON transferencia(idcontacorrente_destino);
CREATE INDEX IF NOT EXISTS idx_tarifa_conta ON tarifa(idcontacorrente);
